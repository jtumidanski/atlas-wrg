package com.atlas.wrg;

import java.net.URI;

import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.UriBuilder;
import com.atlas.wrg.constant.RestConstants;
import com.atlas.wrg.event.consumer.ChannelServerEventConsumer;

public class Server {
   public static void main(String[] args) {
      SimpleEventConsumerFactory.create(new ChannelServerEventConsumer());

      URI uri = UriBuilder.host(RestConstants.SERVICE).uri();
      RestServerFactory.create(uri, "com.atlas.wrg.rest");
   }
}
