package com.atlas.wrg.processor;

import com.atlas.csrv.rest.attribute.ChannelLoadAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataBody;
import rest.DataContainer;

public final class ChannelServerProcessor {
   private ChannelServerProcessor() {
   }

   public static int getLoad(int worldId, int channelId) {
      return UriBuilder.service(RestService.CHANNEL)
            .pathParam("worlds", worldId)
            .pathParam("channels", channelId)
            .path("load")
            .getRestClient(ChannelLoadAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(DataBody::getAttributes)
            .map(ChannelLoadAttributes::load)
            .orElse(0);
   }
}
