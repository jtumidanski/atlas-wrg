package com.atlas.wrg.processor;

import java.util.List;
import java.util.stream.Collectors;

import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.model.ChannelServer;
import com.atlas.wrg.model.ServerStatus;

public class WorldProcessor {
   private static final Object lock = new Object();

   private static volatile WorldProcessor instance;

   public static WorldProcessor getInstance() {
      WorldProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new WorldProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public Integer getCapacityStatus(int worldId) {
      List<ChannelServer> channelServers = ChannelServerRegistry.getInstance().getChannelServers().stream()
            .filter(channelServer -> channelServer.worldId() == worldId)
            .collect(Collectors.toList());
      int channelCount = channelServers.size();

      //TODO create Max players per channel (limit actually used to calculate the World server capacity).
      int max = 100;
      int cap = channelCount * max;
      int count = channelServers.stream()
            .mapToInt(channelServer -> ChannelServerProcessor.getInstance().getLoad(worldId, channelServer.channelId()))
            .sum();

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
}
