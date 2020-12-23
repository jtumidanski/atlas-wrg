package com.atlas.wrg.processor;

import java.util.concurrent.CompletableFuture;
import java.util.stream.Stream;

import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.model.ServerStatus;

public final class WorldProcessor {
   private WorldProcessor() {
   }

   public static Integer getCapacityStatus(int worldId) {
      int channelCount = (int) getChannelServersForWorld(worldId).count();

      //TODO create Max players per channel (limit actually used to calculate the World server capacity).
      int max = 100;
      int cap = channelCount * max;

      int count = totalWorldLoad(worldId);

      ServerStatus serverStatus;
      if (count >= cap) {
         serverStatus = ServerStatus.FULL;
      } else if (count >= cap * 0.8) {
         serverStatus = ServerStatus.HIGHLY_POPULATED;
      } else {
         serverStatus = ServerStatus.NORMAL;
      }
      return serverStatus.getValue();
   }

   protected static Integer totalWorldLoad(int worldId) {
      return getChannelLoadForWorld(worldId)
            .map(CompletableFuture::join)
            .mapToInt(i -> i)
            .sum();
   }

   protected static Stream<CompletableFuture<Integer>> getChannelLoadForWorld(int worldId) {
      return getChannelServersForWorld(worldId)
            .map(ChannelServer::channelId)
            .map(channelId -> ChannelServerProcessor.getLoad(worldId, channelId));
   }

   protected static Stream<ChannelServer> getChannelServersForWorld(int worldId) {
      return ChannelServerRegistry.getInstance().getChannelServers()
            .filter(channelServer -> channelServer.worldId() == worldId);
   }
}
