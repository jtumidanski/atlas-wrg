package com.atlas.wrg.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.wrg.rest.attribute.WorldAttributes;

import builder.AttributeResultBuilder;

public class WorldAttributesBuilder extends RecordBuilder<WorldAttributes, WorldAttributesBuilder>
      implements AttributeResultBuilder {
   private String name;

   private int flag;

   private String message;

   private String eventMessage;

   private boolean recommended;

   private String recommendedMessage;

   private int capacityStatus;

   @Override
   public WorldAttributes construct() {
      return new WorldAttributes(name, flag, message, eventMessage, recommended, recommendedMessage, capacityStatus);
   }

   @Override
   public WorldAttributesBuilder getThis() {
      return this;
   }

   public WorldAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }

   public WorldAttributesBuilder setFlag(int flag) {
      this.flag = flag;
      return getThis();
   }

   public WorldAttributesBuilder setMessage(String message) {
      this.message = message;
      return getThis();
   }

   public WorldAttributesBuilder setEventMessage(String eventMessage) {
      this.eventMessage = eventMessage;
      return getThis();
   }

   public WorldAttributesBuilder setRecommended(boolean recommended) {
      this.recommended = recommended;
      return getThis();
   }

   public WorldAttributesBuilder setRecommendedMessage(String recommendedMessage) {
      this.recommendedMessage = recommendedMessage;
      return getThis();
   }

   public WorldAttributesBuilder setCapacityStatus(int capacityStatus) {
      this.capacityStatus = capacityStatus;
      return getThis();
   }
}
