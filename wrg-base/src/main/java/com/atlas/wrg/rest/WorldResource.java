package com.atlas.wrg.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.model.WorldFlags;
import com.atlas.wrg.processor.ConfigurationProcessor;
import com.atlas.wrg.rest.attribute.WorldAttributes;
import com.atlas.wrg.rest.builder.WorldAttributesBuilder;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;

@Path("worlds")
public class WorldResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getWorldInformation() {

      ResultBuilder resultBuilder = new ResultBuilder();
      ChannelServerRegistry.getInstance().getChannelServers().stream()
            .map(ChannelServer::worldId)
            .distinct()
            .forEach(worldId ->
                  ConfigurationProcessor.getInstance()
                        .getWorldConfiguration(worldId)
                        .ifPresent(configuration -> {
                           WorldFlags worldFlags;
                           try {
                              worldFlags = WorldFlags.valueOf(configuration.flag.toUpperCase());
                           } catch (IllegalArgumentException exception) {
                              System.out
                                    .println("Unable to process world flag configuration for world " + worldId + " "
                                          + "defaulting to Nothing");
                              worldFlags = WorldFlags.NOTHING;
                           }

                           WorldAttributesBuilder builder = new WorldAttributesBuilder()
                                 .setName(configuration.name)
                                 .setFlag(worldFlags.getValue())
                                 .setMessage(configuration.serverMessage)
                                 .setEventMessage(configuration.eventMessage)
                                 .setRecommended(!configuration.whyAmIRecommended.equals(""))
                                 .setRecommendedMessage(configuration.whyAmIRecommended);
                           resultBuilder.addData(new ResultObjectBuilder(WorldAttributes.class, worldId).setAttribute(builder));
                        }));
      return resultBuilder.build();
   }
}
