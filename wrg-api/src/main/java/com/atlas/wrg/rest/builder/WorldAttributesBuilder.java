package com.atlas.wrg.rest.builder;

import com.atlas.wrg.rest.attribute.WorldAttributes;

import builder.AttributeResultBuilder;
import builder.Builder;

public class WorldAttributesBuilder extends Builder<WorldAttributes, WorldAttributesBuilder> implements AttributeResultBuilder {
   @Override
   public WorldAttributes construct() {
      return new WorldAttributes();
   }

   @Override
   public WorldAttributesBuilder getThis() {
      return this;
   }

   public WorldAttributesBuilder setName(String name) {
      return add(attr -> attr.setName(name));
   }

   public WorldAttributesBuilder setFlag(Integer flag) {
      return add(attr -> attr.setFlag(flag));
   }

   public WorldAttributesBuilder setMessage(String message) {
      return add(attr -> attr.setMessage(message));
   }

   public WorldAttributesBuilder setEventMessage(String eventMessage) {
      return add(attr -> attr.setEventMessage(eventMessage));
   }

   public WorldAttributesBuilder setRecommended(Boolean recommended) {
      return add(attr -> attr.setRecommended(recommended));
   }

   public WorldAttributesBuilder setRecommendedMessage(String recommendedMessage) {
      return add(attr -> attr.setRecommendedMessage(recommendedMessage));
   }

   public WorldAttributesBuilder setCapacityStatus(Integer capacityStatus) {
      return add(attr -> attr.setCapacityStatus(capacityStatus));
   }
}
