package com.atlas.wrg.event.consumer;

import com.atlas.csrv.event.ChannelServerEvent;
import com.atlas.csrv.event.ChannelServerEventStatus;
import com.atlas.kafka.consumer.ConsumerRecordHandler;
import com.atlas.wrg.ChannelServerRegistry;

public class ChannelServerEventConsumer implements ConsumerRecordHandler<Long, ChannelServerEvent> {
   @Override
   public void handle(Long key, ChannelServerEvent event) {
      if (event.status() == ChannelServerEventStatus.STARTED) {
         ChannelServerRegistry.getInstance().addChannelServer(event.worldId(), event.channelId(), event.ipAddress(), event.port());
      } else if (event.status() == ChannelServerEventStatus.SHUTDOWN) {
         ChannelServerRegistry.getInstance().removeChannelServer(event.worldId(), event.channelId());
      }
   }
}
