package com.atlas.wrg.rest;

import java.util.List;
import java.util.Optional;
import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;
import com.atlas.wrg.rest.builder.ChannelServerAttributesBuilder;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;
import rest.InputBody;

@Path("channelServers")
public class ChannelServerResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getRegisteredChannelServers() {
      List<ChannelServer> channelServers = ChannelServerRegistry.getInstance().getChannelServers();
      ResultBuilder resultBuilder = new ResultBuilder();
      channelServers.forEach(channelServer -> resultBuilder.addData(produceResult(channelServer)));
      return resultBuilder.build();
   }

   @POST
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response registerChannelServer(InputBody<ChannelServerAttributes> inputBody) {
      int worldId = inputBody.attribute(ChannelServerAttributes::getWorldId);
      int channelId = inputBody.attribute(ChannelServerAttributes::getChannelId);
      Optional<ChannelServer> channelServer = ChannelServerRegistry.getInstance().addChannelServer(worldId, channelId);
      ResultBuilder resultBuilder = new ResultBuilder(Response.Status.CONFLICT);
      if (channelServer.isPresent()) {
         resultBuilder = new ResultBuilder(Response.Status.CREATED);
         resultBuilder.addData(produceResult(channelServer.get()));
      }
      return resultBuilder.build();
   }

   protected ResultObjectBuilder produceResult(ChannelServer channelServer) {
      return new ResultObjectBuilder(ChannelServerAttributes.class, channelServer.uniqueId())
            .setAttribute(new ChannelServerAttributesBuilder()
                  .setWorldId(channelServer.worldId())
                  .setChannelId(channelServer.channelId()));
   }
}
