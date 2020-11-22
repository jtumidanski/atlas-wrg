package com.atlas.wrg;

import java.time.Duration;
import java.util.Collections;

import com.atlas.csrv.event.ChannelServer;
import com.atlas.csrv.event.ChannelServerStatus;
import com.atlas.kafka.KafkaConsumerFactory;
import com.atlas.wrg.constant.KafkaConstants;
import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerRecords;

public class ChannelServerEventConsumer implements Runnable {
   @Override
   public void run() {
      Consumer<Long, ChannelServer> consumer = KafkaConsumerFactory.createConsumer("World Registry",
            KafkaConstants.BOOTSTRAP_SERVERS, ChannelServer.class);
      consumer.subscribe(Collections.singleton(KafkaConstants.CHANNEL_SERVICE_STATUS_TOPIC));

      try (consumer) {
         while (true) {
            final ConsumerRecords<Long, ChannelServer> consumerRecords = consumer.poll(Duration.ofMillis(1000));

            if (consumerRecords.count() > 0) {
               consumerRecords.forEach(record -> {
                  if (record.value().status() == ChannelServerStatus.STARTED) {
                     ChannelServerRegistry.getInstance().addChannelServer(record.value().worldId(), record.value().channelId(),
                           record.value().ipAddress(), record.value().port());
                  } else if (record.value().status() == ChannelServerStatus.SHUTDOWN) {
                     ChannelServerRegistry.getInstance().removeChannelServer(record.value().worldId(), record.value().channelId());
                  }
               });
               consumer.commitAsync();
            }
         }
      }
   }
}
