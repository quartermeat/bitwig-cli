package dev.codextools.bitwigcli;

import com.bitwig.extension.controller.ControllerExtension;
import com.bitwig.extension.controller.api.Application;
import com.bitwig.extension.controller.api.ClipLauncherSlotBank;
import com.bitwig.extension.controller.api.ControllerHost;
import com.bitwig.extension.controller.api.Project;
import com.bitwig.extension.controller.api.SceneBank;
import com.bitwig.extension.controller.api.Track;
import com.bitwig.extension.controller.api.TrackBank;
import com.bitwig.extension.controller.api.Transport;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.nio.charset.StandardCharsets;

public final class BitwigCliBridgeExtension extends ControllerExtension {
    private static final int PORT = 58770;

    private ControllerHost host;
    private Application application;
    private Project project;
    private Transport transport;
    private TrackBank trackBank;
    private Thread listenerThread;
    private DatagramSocket socket;
    private volatile boolean running;

    protected BitwigCliBridgeExtension(final BitwigCliBridgeDefinition definition, final ControllerHost host) {
        super(definition, host);
    }

    @Override
    public void init() {
        host = getHost();
        application = host.createApplication();
        project = host.getProject();
        transport = host.createTransport();
        trackBank = host.createMainTrackBank(16, 0, 16);
        running = true;
        startListener();
        host.println("Bitwig CLI Bridge listening on UDP 127.0.0.1:" + PORT);
    }

    @Override
    public void exit() {
        running = false;
        if (socket != null) {
            socket.close();
        }
        host.println("Bitwig CLI Bridge stopped");
    }

    @Override
    public void flush() {
    }

    private void startListener() {
        listenerThread = new Thread(this::listen, "bitwig-cli-bridge");
        listenerThread.setDaemon(true);
        listenerThread.start();
    }

    private void listen() {
        try (DatagramSocket server = new DatagramSocket(PORT, InetAddress.getByName("127.0.0.1"))) {
            socket = server;
            while (running) {
                byte[] buffer = new byte[4096];
                DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                server.receive(packet);

                String body = new String(packet.getData(), packet.getOffset(), packet.getLength(), StandardCharsets.UTF_8);
                String reply = handle(body);
                byte[] response = reply.getBytes(StandardCharsets.UTF_8);
                server.send(new DatagramPacket(response, response.length, packet.getAddress(), packet.getPort()));
            }
        } catch (IOException ex) {
            if (running) {
                host.println("Bitwig CLI Bridge listener failed: " + ex.getMessage());
            }
        }
    }

    private String handle(String body) {
        Command command = Command.parse(body);
        host.scheduleTask(() -> apply(command), 0);
        return "queued " + command.action;
    }

    private void apply(Command command) {
        switch (command.action) {
            case "status":
                host.println("Bitwig CLI Bridge status requested");
                break;
            case "play":
                transport.play();
                break;
            case "stop":
            case "panic":
                transport.stop();
                break;
            case "record":
                transport.record();
                break;
            case "overdub":
                transport.isClipLauncherOverdubEnabled().toggle();
                break;
            case "tempo":
                transport.tempo().set(command.value);
                break;
            case "track.arm":
                track(command.track).arm().toggle();
                break;
            case "track.mute":
                track(command.track).mute().toggle();
                break;
            case "track.solo":
                track(command.track).solo().toggle();
                break;
            case "track.select":
                track(command.track).selectInEditor();
                break;
            case "track.create.instrument":
                application.createInstrumentTrack(zeroBasedPosition(command.track));
                break;
            case "track.create.audio":
                application.createAudioTrack(zeroBasedPosition(command.track));
                break;
            case "track.create.effect":
                application.createEffectTrack(zeroBasedPosition(command.track));
                break;
            case "scene.create":
                project.createScene();
                break;
            case "scene.select":
                scenes().getScene(command.slot - 1).selectInEditor();
                break;
            case "scene.show":
                scenes().getScene(command.slot - 1).showInEditor();
                break;
            case "scene.launch":
                scenes().launchScene(command.slot - 1);
                break;
            case "clip.launch":
                slots(command.track).getItemAt(command.slot - 1).launch();
                break;
            case "clip.record":
                slots(command.track).getItemAt(command.slot - 1).record();
                break;
            case "clip.stop":
                track(command.track).stop();
                break;
            case "clip.select":
                slots(command.track).getItemAt(command.slot - 1).select();
                break;
            case "clip.show":
                slots(command.track).getItemAt(command.slot - 1).showInEditor();
                break;
            case "clip.duplicate":
                slots(command.track).getItemAt(command.slot - 1).duplicateClip();
                break;
            case "clip.delete":
                slots(command.track).deleteClip(command.slot - 1);
                break;
            case "clip.create":
                slots(command.track).createEmptyClip(command.slot - 1, (int) command.value);
                slots(command.track).getItemAt(command.slot - 1).select();
                slots(command.track).getItemAt(command.slot - 1).showInEditor();
                break;
            default:
                host.println("Unknown Bitwig CLI command: " + command.action);
                break;
        }
    }

    private Track track(int oneBasedTrack) {
        return trackBank.getItemAt(oneBasedTrack - 1);
    }

    private ClipLauncherSlotBank slots(int oneBasedTrack) {
        return track(oneBasedTrack).clipLauncherSlotBank();
    }

    private SceneBank scenes() {
        return trackBank.sceneBank();
    }

    private int zeroBasedPosition(int oneBasedPosition) {
        return oneBasedPosition > 0 ? oneBasedPosition - 1 : -1;
    }

    private static final class Command {
        String action = "";
        int track;
        int slot;
        double value;

        static Command parse(String body) {
            Command command = new Command();
            command.action = extractString(body, "action");
            command.track = extractInt(body, "track");
            command.slot = extractInt(body, "slot");
            command.value = extractDouble(body, "value");
            return command;
        }

        private static String extractString(String body, String key) {
            String marker = "\"" + key + "\":\"";
            int start = body.indexOf(marker);
            if (start < 0) {
                return "";
            }
            start += marker.length();
            int end = body.indexOf('"', start);
            return end < 0 ? "" : body.substring(start, end);
        }

        private static int extractInt(String body, String key) {
            return (int) extractDouble(body, key);
        }

        private static double extractDouble(String body, String key) {
            String marker = "\"" + key + "\":";
            int start = body.indexOf(marker);
            if (start < 0) {
                return 0;
            }
            start += marker.length();
            int end = start;
            while (end < body.length() && "-0123456789.".indexOf(body.charAt(end)) >= 0) {
                end++;
            }
            try {
                return Double.parseDouble(body.substring(start, end));
            } catch (NumberFormatException ex) {
                return 0;
            }
        }
    }
}
