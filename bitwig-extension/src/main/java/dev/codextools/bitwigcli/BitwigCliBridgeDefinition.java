package dev.codextools.bitwigcli;

import com.bitwig.extension.api.PlatformType;
import com.bitwig.extension.controller.AutoDetectionMidiPortNamesList;
import com.bitwig.extension.controller.ControllerExtension;
import com.bitwig.extension.controller.ControllerExtensionDefinition;
import com.bitwig.extension.controller.api.ControllerHost;

import java.util.UUID;

public final class BitwigCliBridgeDefinition extends ControllerExtensionDefinition {
    private static final UUID ID = UUID.fromString("655fc4fe-9a82-4ccb-9ff2-d3f690bdeedb");

    @Override
    public String getName() {
        return "Bitwig CLI Bridge";
    }

    @Override
    public String getAuthor() {
        return "Codex Tools";
    }

    @Override
    public String getVersion() {
        return "0.1.0";
    }

    @Override
    public UUID getId() {
        return ID;
    }

    @Override
    public String getHardwareVendor() {
        return "Codex Tools";
    }

    @Override
    public String getHardwareModel() {
        return "Bitwig CLI Bridge";
    }

    @Override
    public int getRequiredAPIVersion() {
        return 18;
    }

    @Override
    public int getNumMidiInPorts() {
        return 0;
    }

    @Override
    public int getNumMidiOutPorts() {
        return 0;
    }

    @Override
    public void listAutoDetectionMidiPortNames(final AutoDetectionMidiPortNamesList list, final PlatformType platformType) {
    }

    @Override
    public ControllerExtension createInstance(ControllerHost host) {
        return new BitwigCliBridgeExtension(this, host);
    }
}
