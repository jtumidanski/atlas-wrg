package com.atlas.wrg;

import java.net.URI;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

public class Server {

   private static final ExecutorService executorService = Executors.newSingleThreadExecutor();

   public static void main(String[] args) {
      executorService.execute(new ChannelServerEventConsumer());

      URI uri = UriBuilder.host(RestService.WORLD_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.wrg.rest");
   }
}
