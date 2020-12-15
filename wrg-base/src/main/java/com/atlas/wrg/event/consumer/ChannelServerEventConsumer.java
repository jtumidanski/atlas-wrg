package com.atlas.wrg.event.consumer;

import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.ChannelServerEvent;
import com.atlas.csrv.event.ChannelServerEventStatus;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.wrg.ChannelServerRegistry;
import com.atlas.wrg.processor.TopicDiscoveryProcessor;

public class ChannelServerEventConsumer implements SimpleEventHandler<ChannelServerEvent> {
   @Override
   public void handle(Long key, ChannelServerEvent event) {
      if (event.status() == ChannelServerEventStatus.STARTED) {
         ChannelServerRegistry.getInstance().addChannelServer(event.worldId(), event.channelId(), event.ipAddress(), event.port());
      } else if (event.status() == ChannelServerEventStatus.SHUTDOWN) {
         ChannelServerRegistry.getInstance().removeChannelServer(event.worldId(), event.channelId());
      }
   }

   @Override
   public Class<ChannelServerEvent> getEventClass() {
      return ChannelServerEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "World Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_CHANNEL_SERVICE);
   }
}
