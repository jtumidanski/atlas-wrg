package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public record ChannelServerAttributes(Integer worldId, Integer channelId, Integer capacity, String ipAddress, Integer port)
      implements AttributeResult {
}
