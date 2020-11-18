package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public record ChannelServerAttributes(int worldId, int channelId, int capacity) implements AttributeResult {
}
