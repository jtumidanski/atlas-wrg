package com.atlas.wrg;

import java.net.URI;

import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

public class Server {
   public static void main(String[] args) {
      URI uri = UriBuilder.host(RestService.WORLD_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.wrg.rest");
   }
}
