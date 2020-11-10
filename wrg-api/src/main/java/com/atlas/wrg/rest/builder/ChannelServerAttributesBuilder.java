package com.atlas.wrg.rest.builder;

import com.atlas.wrg.rest.attribute.ChannelServerAttributes;

import builder.AttributeResultBuilder;
import builder.Builder;

public class ChannelServerAttributesBuilder extends Builder<ChannelServerAttributes, ChannelServerAttributesBuilder>
      implements AttributeResultBuilder {
   @Override
   public ChannelServerAttributes construct() {
      return new ChannelServerAttributes();
   }

   @Override
   public ChannelServerAttributesBuilder getThis() {
      return this;
   }

   public ChannelServerAttributesBuilder setWorldId(Integer worldId) {
      return add(attr -> attr.setWorldId(worldId));
   }

   public ChannelServerAttributesBuilder setChannelId(Integer channelId) {
      return add(attr -> attr.setChannelId(channelId));
   }
}
