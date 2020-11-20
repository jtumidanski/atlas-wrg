package com.atlas.wrg.rest;

import java.util.Optional;

import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.model.WorldFlags;
import com.atlas.wrg.processor.ChannelServerProcessor;
import com.atlas.wrg.processor.ConfigurationProcessor;
import com.atlas.wrg.processor.WorldProcessor;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;
import com.atlas.wrg.rest.attribute.WorldAttributes;
import com.atlas.wrg.rest.builder.ChannelServerAttributesBuilder;
import com.atlas.wrg.rest.builder.WorldAttributesBuilder;

import builder.ResultObjectBuilder;

public final class ResultObjectFactory {
   public static ResultObjectBuilder create(ChannelServer channelServer) {
      ChannelServerAttributesBuilder builder = new ChannelServerAttributesBuilder()
            .setWorldId(channelServer.worldId())
            .setChannelId(channelServer.channelId())
            .setIpAddress(channelServer.ipAddress())
            .setPort(channelServer.port());

      int channelLoad = ChannelServerProcessor.getInstance().getLoad(channelServer.worldId(), channelServer.channelId());
      builder.setCapacity(channelLoad);

      return new ResultObjectBuilder(ChannelServerAttributes.class, channelServer.uniqueId())
            .setAttribute(builder);
   }

   public static Optional<ResultObjectBuilder> create(Integer worldId) {
      return ConfigurationProcessor.getInstance()
            .getWorldConfiguration(worldId)
            .map(configuration -> {
               WorldFlags worldFlags;
               try {
                  worldFlags = WorldFlags.valueOf(configuration.flag.toUpperCase());
               } catch (IllegalArgumentException exception) {
                  System.out
                        .println("Unable to process world flag configuration for world " + worldId + " "
                              + "defaulting to Nothing");
                  worldFlags = WorldFlags.NOTHING;
               }

               Integer capacityStatus = WorldProcessor.getInstance().getCapacityStatus(worldId);

               WorldAttributesBuilder builder = new WorldAttributesBuilder()
                     .setName(configuration.name)
                     .setFlag(worldFlags.getValue())
                     .setMessage(configuration.serverMessage)
                     .setEventMessage(configuration.eventMessage)
                     .setRecommended(!configuration.whyAmIRecommended.equals(""))
                     .setRecommendedMessage(configuration.whyAmIRecommended)
                     .setCapacityStatus(capacityStatus);
               return new ResultObjectBuilder(WorldAttributes.class, worldId).setAttribute(builder);
            });
   }
}
