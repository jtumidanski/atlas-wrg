package com.atlas.wrg.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;

import builder.AttributeResultBuilder;

public class ChannelServerAttributesBuilder extends RecordBuilder<ChannelServerAttributes, ChannelServerAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer worldId;

   private Integer channelId;

   private Integer capacity;

   @Override
   public ChannelServerAttributes construct() {
      return new ChannelServerAttributes(worldId, channelId, capacity);
   }

   @Override
   public ChannelServerAttributesBuilder getThis() {
      return this;
   }

   public ChannelServerAttributesBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public ChannelServerAttributesBuilder setChannelId(Integer channelId) {
      this.channelId = channelId;
      return getThis();
   }

   public ChannelServerAttributesBuilder setCapacity(Integer capacity) {
      this.capacity = capacity;
      return getThis();
   }
}
