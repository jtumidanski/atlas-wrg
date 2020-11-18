package com.atlas.wrg.rest.builder;

import builder.AttributeResultBuilder;
import builder.RecordBuilder;
import com.atlas.wrg.rest.attribute.WorldAttributes;

public class WorldAttributesBuilder extends RecordBuilder<WorldAttributes, WorldAttributesBuilder> implements AttributeResultBuilder {
   private static final String NAME = "NAME";

   private static final String FLAG = "FLAG";

   private static final String MESSAGE = "MESSAGE";

   private static final String EVENT_MESSAGE = "EVENT_MESSAGE";

   private static final String RECOMMENDED = "RECOMMENDED";

   private static final String RECOMMENDED_MESSAGE = "RECOMMENDED_MESSAGE";

   private static final String CAPACITY_STATUS = "CAPACITY_STATUS";


   @Override
   public WorldAttributes construct() {
      return new WorldAttributes(get(NAME), get(FLAG), get(MESSAGE), get(EVENT_MESSAGE), get(RECOMMENDED), get(RECOMMENDED_MESSAGE), get(CAPACITY_STATUS));
   }

   @Override
   public WorldAttributesBuilder getThis() {
      return this;
   }

   public WorldAttributesBuilder setName(String name) {
      return set(NAME, name);
   }

   public WorldAttributesBuilder setFlag(int flag) {
      return set(FLAG, flag);
   }

   public WorldAttributesBuilder setMessage(String message) {
      return set(MESSAGE, message);
   }

   public WorldAttributesBuilder setEventMessage(String eventMessage) {
      return set(EVENT_MESSAGE, eventMessage);
   }

   public WorldAttributesBuilder setRecommended(boolean recommended) {
      return set(RECOMMENDED, recommended);
   }

   public WorldAttributesBuilder setRecommendedMessage(String recommendedMessage) {
      return set(RECOMMENDED_MESSAGE, recommendedMessage);
   }

   public WorldAttributesBuilder setCapacityStatus(int capacityStatus) {
      return set(CAPACITY_STATUS, capacityStatus);
   }
}
