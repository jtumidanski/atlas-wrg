package com.atlas.wrg.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;

import builder.AttributeResultBuilder;

public class ChannelServerAttributesBuilder extends RecordBuilder<ChannelServerAttributes, ChannelServerAttributesBuilder>
      implements AttributeResultBuilder {
   private int worldId;

   private int channelId;

   private int capacity;

   @Override
   public ChannelServerAttributes construct() {
      return new ChannelServerAttributes(worldId, channelId, capacity);
   }

   @Override
   public ChannelServerAttributesBuilder getThis() {
      return this;
   }

   public ChannelServerAttributesBuilder setWorldId(int worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public ChannelServerAttributesBuilder setChannelId(int channelId) {
      this.channelId = channelId;
      return getThis();
   }

   public ChannelServerAttributesBuilder setCapacity(int capacity) {
      this.capacity = capacity;
      return getThis();
   }
}
