package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public class ChannelServerAttributes implements AttributeResult {
   private Integer worldId;

   private String worldName;

   private String flag;

   private String eventMessage;

   private Integer channelId;

   private Integer capacity;

   public Integer getWorldId() {
      return worldId;
   }

   public void setWorldId(Integer worldId) {
      this.worldId = worldId;
   }

   public String getWorldName() {
      return worldName;
   }

   public void setWorldName(String worldName) {
      this.worldName = worldName;
   }

   public String getFlag() {
      return flag;
   }

   public void setFlag(String flag) {
      this.flag = flag;
   }

   public String getEventMessage() {
      return eventMessage;
   }

   public void setEventMessage(String eventMessage) {
      this.eventMessage = eventMessage;
   }

   public Integer getChannelId() {
      return channelId;
   }

   public void setChannelId(Integer channelId) {
      this.channelId = channelId;
   }

   public Integer getCapacity() {
      return capacity;
   }

   public void setCapacity(Integer capacity) {
      this.capacity = capacity;
   }
}
