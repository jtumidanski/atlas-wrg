package com.atlas.wrg.model;

public record ChannelServer(int uniqueId, int worldId, int channelId, String ipAddress, int port) {
}
