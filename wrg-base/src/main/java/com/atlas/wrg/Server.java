package com.atlas.wrg;

import java.net.URI;
import java.util.concurrent.Executors;

import com.atlas.csrv.event.ChannelServerEvent;
import com.atlas.kafka.consumer.ConsumerBuilder;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

public class Server {
   public static void main(String[] args) {
      Executors.newSingleThreadExecutor().execute(new ConsumerBuilder<>("World Registry", ChannelServerEvent.class)
            .setBootstrapServers(System.getenv("BOOTSTRAP_SERVERS"))
            .setTopic(System.getenv("TOPIC_CHANNEL_SERVICE"))
            .setHandler(new com.atlas.wrg.event.consumer.ChannelServerEventConsumer())
            .build()
      );

      URI uri = UriBuilder.host(RestService.WORLD_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.wrg.rest");
   }
}
