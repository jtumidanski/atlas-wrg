package com.atlas.wrg.rest.builder;

import builder.AttributeResultBuilder;
import builder.RecordBuilder;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;

public class ChannelServerAttributesBuilder extends RecordBuilder<ChannelServerAttributes, ChannelServerAttributesBuilder> implements AttributeResultBuilder {
   private static final String WORLD_ID = "WORLD_ID";

   private static final String CHANNEL_ID = "CHANNEL_ID";

   private static final String CAPACITY = "CAPACITY";


   @Override
   public ChannelServerAttributes construct() {
      return new ChannelServerAttributes(get(WORLD_ID), get(CHANNEL_ID), get(CAPACITY));
   }

   @Override
   public ChannelServerAttributesBuilder getThis() {
      return this;
   }

   public ChannelServerAttributesBuilder setWorldId(int worldId) {
      return set(WORLD_ID, worldId);
   }

   public ChannelServerAttributesBuilder setChannelId(int channelId) {
      return set(CHANNEL_ID, channelId);
   }

   public ChannelServerAttributesBuilder setCapacity(int capacity) {
      return set(CAPACITY, capacity);
   }
}
