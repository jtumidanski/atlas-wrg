package com.atlas.wrg;

import java.time.Duration;
import java.util.Collections;

import com.atlas.csrv.event.ChannelServerEvent;
import com.atlas.csrv.event.ChannelServerEventStatus;
import com.atlas.kafka.KafkaConsumerFactory;
import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerRecords;

public class ChannelServerEventConsumer implements Runnable {
   private final String channelServiceTopic;

   private final String bootstrapServers;

   public ChannelServerEventConsumer(String bootstrapServers, String channelServiceTopic) {
      this.bootstrapServers = bootstrapServers;
      this.channelServiceTopic = channelServiceTopic;
   }

   @Override
   public void run() {
      Consumer<Long, ChannelServerEvent> consumer = KafkaConsumerFactory.createConsumer("World Registry",
            bootstrapServers, ChannelServerEvent.class);
      consumer.subscribe(Collections.singleton(channelServiceTopic));

      try (consumer) {
         while (true) {
            final ConsumerRecords<Long, ChannelServerEvent> consumerRecords = consumer.poll(Duration.ofMillis(1000));

            if (consumerRecords.count() > 0) {
               consumerRecords.forEach(record -> {
                  if (record.value().status() == ChannelServerEventStatus.STARTED) {
                     ChannelServerRegistry.getInstance().addChannelServer(record.value().worldId(), record.value().channelId(),
                           record.value().ipAddress(), record.value().port());
                  } else if (record.value().status() == ChannelServerEventStatus.SHUTDOWN) {
                     ChannelServerRegistry.getInstance().removeChannelServer(record.value().worldId(), record.value().channelId());
                  }
               });
               consumer.commitAsync();
            }
         }
      }
   }
}
