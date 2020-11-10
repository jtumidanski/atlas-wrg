package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public class ChannelServerAttributes implements AttributeResult {
   private Integer worldId;

   private Integer channelId;

   public Integer getWorldId() {
      return worldId;
   }

   public void setWorldId(Integer worldId) {
      this.worldId = worldId;
   }

   public Integer getChannelId() {
      return channelId;
   }

   public void setChannelId(Integer channelId) {
      this.channelId = channelId;
   }
}
