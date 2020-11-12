package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public class WorldAttributes implements AttributeResult {
   private String name;

   private Integer flag;

   private String message;

   private String eventMessage;

   private Boolean recommended;

   private String recommendedMessage;

   private Integer capacityStatus;

   public String getName() {
      return name;
   }

   public void setName(String name) {
      this.name = name;
   }

   public Integer getFlag() {
      return flag;
   }

   public void setFlag(Integer flag) {
      this.flag = flag;
   }

   public String getMessage() {
      return message;
   }

   public void setMessage(String message) {
      this.message = message;
   }

   public String getEventMessage() {
      return eventMessage;
   }

   public void setEventMessage(String eventMessage) {
      this.eventMessage = eventMessage;
   }

   public Boolean getRecommended() {
      return recommended;
   }

   public void setRecommended(Boolean recommended) {
      this.recommended = recommended;
   }

   public String getRecommendedMessage() {
      return recommendedMessage;
   }

   public void setRecommendedMessage(String recommendedMessage) {
      this.recommendedMessage = recommendedMessage;
   }

   public Integer getCapacityStatus() {
      return capacityStatus;
   }

   public void setCapacityStatus(Integer capacityStatus) {
      this.capacityStatus = capacityStatus;
   }
}
