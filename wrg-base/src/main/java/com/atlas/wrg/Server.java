package com.atlas.wrg;

import java.net.URI;

import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;
import com.atlas.wrg.event.consumer.ChannelServerEventConsumer;

public class Server {
   public static void main(String[] args) {
      SimpleEventConsumerFactory.create(new ChannelServerEventConsumer());

      URI uri = UriBuilder.host(RestService.WORLD_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.wrg.rest");
   }
}
