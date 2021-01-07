package com.atlas.wrg.processor;

import java.util.concurrent.CompletableFuture;

import com.app.rest.util.RestResponseUtil;
import com.atlas.csrv.constant.RestConstants;
import com.atlas.csrv.rest.attribute.ChannelLoadAttributes;
import com.atlas.shared.rest.UriBuilder;

import rest.DataBody;
import rest.DataContainer;

public final class ChannelServerProcessor {
   private ChannelServerProcessor() {
   }

   protected static CompletableFuture<DataContainer<ChannelLoadAttributes>> requestChannelLoad(int worldId, int channelId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("worlds", worldId)
            .pathParam("channels", channelId)
            .path("load")
            .getAsyncRestClient(ChannelLoadAttributes.class)
            .get()
            .thenApply(RestResponseUtil::result);
   }

   protected static Integer getLoadFromContainer(DataContainer<ChannelLoadAttributes> container) {
      return container.data()
            .map(DataBody::getAttributes)
            .map(ChannelLoadAttributes::load)
            .orElse(0);
   }

   public static CompletableFuture<Integer> getLoad(int worldId, int channelId) {
      return requestChannelLoad(worldId, channelId)
            .thenApply(ChannelServerProcessor::getLoadFromContainer);
   }
}
