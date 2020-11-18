package com.atlas.wrg.processor;

import com.atlas.csrv.attribute.ChannelLoadAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataBody;
import rest.DataContainer;

public class ChannelServerProcessor {
   private static final Object lock = new Object();

   private static volatile ChannelServerProcessor instance;

   public static ChannelServerProcessor getInstance() {
      ChannelServerProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ChannelServerProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public int getLoad(int worldId, int channelId) {
      return UriBuilder.service(RestService.CHANNEL)
            .path("worlds").path(worldId)
            .path("channels").path(channelId)
            .path("load")
            .getRestClient(ChannelLoadAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getData)
            .map(DataBody::getAttributes)
            .map(ChannelLoadAttributes::load)
            .orElse(0);
   }
}
