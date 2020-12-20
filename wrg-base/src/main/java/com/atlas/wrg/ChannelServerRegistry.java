package com.atlas.wrg;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import com.atlas.wrg.model.ChannelServer;

public class ChannelServerRegistry {
   private static final Object lock = new Object();

   private static volatile ChannelServerRegistry instance;

   private static final Object registryLock = new Object();

   private final AtomicInteger runningUniqueId = new AtomicInteger(1000000001);

   private final List<ChannelServer> channelServerList = new ArrayList<>();

   public static ChannelServerRegistry getInstance() {
      ChannelServerRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ChannelServerRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   public Optional<ChannelServer> addChannelServer(int worldId, int channelId, String ipAddress, int port) {
      synchronized (registryLock) {
         boolean alreadyExists = channelServerList.stream()
               .anyMatch(server -> server.worldId() == worldId && server.channelId() == channelId);
         if (alreadyExists) {
            return Optional.empty();
         }

         List<Integer> existingIds = channelServerList.stream()
               .map(ChannelServer::uniqueId)
               .collect(Collectors.toList());

         int currentUniqueId;
         do {
            if ((currentUniqueId = runningUniqueId.incrementAndGet()) >= 2000000000) {
               runningUniqueId.set(currentUniqueId = 1000000001);
            }
         } while (existingIds.contains(currentUniqueId));

         ChannelServer channelServer = new ChannelServer(currentUniqueId, worldId, channelId, ipAddress, port);
         channelServerList.add(channelServer);
         return Optional.of(channelServer);
      }
   }

   public Stream<ChannelServer> getChannelServers() {
      return Collections.unmodifiableList(channelServerList).stream();
   }

   public void removeChannelServer(Integer id) {
      synchronized (registryLock) {
         channelServerList.stream()
               .filter(possible -> possible.uniqueId() == id)
               .findFirst()
               .ifPresent(channelServerList::remove);
      }
   }

   public void removeChannelServer(int worldId, int channelId) {
      synchronized (registryLock) {
         channelServerList.stream()
               .filter(possible -> possible.worldId() == worldId)
               .filter(possible -> possible.channelId() == channelId)
               .findFirst()
               .ifPresent(channelServerList::remove);
      }
   }

   public Optional<ChannelServer> getChannelServer(int worldId, int channelId) {
      return getChannelServers()
            .filter(server -> server.worldId() == worldId)
            .filter(server -> server.channelId() == channelId)
            .findFirst();
   }
}
