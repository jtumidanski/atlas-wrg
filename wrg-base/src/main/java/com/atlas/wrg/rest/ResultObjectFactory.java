package com.atlas.wrg.rest;

import builder.ResultObjectBuilder;
import com.atlas.wrg.ConfigurationRegistry;
import com.atlas.wrg.configuration.WorldConfiguration;
import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.model.WorldFlags;
import com.atlas.wrg.processor.ChannelServerProcessor;
import com.atlas.wrg.processor.WorldProcessor;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;
import com.atlas.wrg.rest.attribute.WorldAttributes;
import com.atlas.wrg.rest.builder.ChannelServerAttributesBuilder;
import com.atlas.wrg.rest.builder.WorldAttributesBuilder;

import java.util.Optional;

public final class ResultObjectFactory {
   public static ResultObjectBuilder create(ChannelServer channelServer) {
      int channelLoad = ChannelServerProcessor.getLoad(channelServer.worldId(), channelServer.channelId());
      return new ResultObjectBuilder(ChannelServerAttributes.class, channelServer.uniqueId())
            .setAttribute(new ChannelServerAttributesBuilder()
                  .setWorldId(channelServer.worldId())
                  .setChannelId(channelServer.channelId())
                  .setIpAddress(channelServer.ipAddress())
                  .setPort(channelServer.port())
                  .setCapacity(channelLoad)
            );
   }

   public static Optional<ResultObjectBuilder> create(Integer worldId) {
      return ConfigurationRegistry.getInstance()
            .getWorldConfiguration(worldId)
            .map(configuration -> createWorldFromConfiguration(worldId, configuration));
   }

   private static ResultObjectBuilder createWorldFromConfiguration(Integer worldId, WorldConfiguration configuration) {
      WorldFlags worldFlags;
      try {
         worldFlags = WorldFlags.valueOf(configuration.flag.toUpperCase());
      } catch (IllegalArgumentException exception) {
         System.out
               .println("Unable to process world flag configuration for world " + worldId + " "
                     + "defaulting to Nothing");
         worldFlags = WorldFlags.NOTHING;
      }

      Integer capacityStatus = WorldProcessor.getCapacityStatus(worldId);

      return new ResultObjectBuilder(WorldAttributes.class, worldId)
            .setAttribute(new WorldAttributesBuilder()
                  .setName(configuration.name)
                  .setFlag(worldFlags.getValue())
                  .setMessage(configuration.serverMessage)
                  .setEventMessage(configuration.eventMessage)
                  .setRecommended(!configuration.whyAmIRecommended.equals(""))
                  .setRecommendedMessage(configuration.whyAmIRecommended)
                  .setCapacityStatus(capacityStatus)
            );
   }
}
