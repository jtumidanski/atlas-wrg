package com.atlas.wrg.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.rest.attribute.ChannelServerAttributes;

import builder.ResultBuilder;
import rest.InputBody;

@Path("channelServers")
public class ChannelServerResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getRegisteredChannelServers(@QueryParam("world") Integer worldId) {
      if (worldId != null) {
         return ChannelServerRegistry.getInstance().getChannelServers().stream()
               .filter(channelServer -> channelServer.worldId() == worldId)
               .map(ResultObjectFactory::create)
               .collect(Collectors.toResultBuilder())
               .build();
      } else {
         return ChannelServerRegistry.getInstance().getChannelServers().stream()
               .map(ResultObjectFactory::create)
               .collect(Collectors.toResultBuilder())
               .build();
      }
   }

   @POST
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response registerChannelServer(InputBody<ChannelServerAttributes> inputBody) {
      int worldId = inputBody.attribute(ChannelServerAttributes::worldId);
      int channelId = inputBody.attribute(ChannelServerAttributes::channelId);
      String ipAddress = inputBody.attribute(ChannelServerAttributes::ipAddress);
      int port = inputBody.attribute(ChannelServerAttributes::port);

      return ChannelServerRegistry.getInstance()
            .addChannelServer(worldId, channelId, ipAddress, port)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.CONFLICT))
            .build();
   }

   @DELETE
   @Path("/{id}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response registerChannelServer(@PathParam("id") Integer id) {
      ChannelServerRegistry.getInstance().removeChannelServer(id);
      return new ResultBuilder(Response.Status.NO_CONTENT).build();
   }
}
