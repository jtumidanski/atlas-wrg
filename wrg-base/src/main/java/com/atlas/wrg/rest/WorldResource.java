package com.atlas.wrg.rest;

import java.util.Optional;
import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.model.ChannelServer;

import builder.ResultBuilder;

@Path("worlds")
public class WorldResource {
   @GET
   @Path("")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getWorldInformation() {
      return ChannelServerRegistry.getInstance().getChannelServers().stream()
            .map(ChannelServer::worldId)
            .distinct()
            .map(ResultObjectFactory::create)
            .flatMap(Optional::stream)
            .collect(Collectors.toResultBuilder())
            .build();
   }

   @GET
   @Path("/{worldId}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getWorldInformation(@PathParam("worldId") Integer worldId) {
      return ResultObjectFactory.create(worldId)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND))
            .build();
   }
}
