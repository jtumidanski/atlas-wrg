package com.atlas.wrg.processor;

import java.util.Optional;

import com.app.rest.util.RestResponseUtil;
import com.atlas.shared.rest.UriBuilder;
import com.atlas.tds.constant.RestConstants;
import com.atlas.tds.rest.attribute.TopicAttributes;

import rest.DataBody;
import rest.DataContainer;

public final class TopicDiscoveryProcessor {
   private TopicDiscoveryProcessor() {
   }

   public static String getTopic(String id) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("topics", id)
            .getAsyncRestClient(TopicAttributes.class)
            .retryOnFailure(1000)
            .get()
            .thenApply(RestResponseUtil::result)
            .thenApply(DataContainer::data)
            .thenApply(Optional::get)
            .thenApply(DataBody::getAttributes)
            .thenApply(TopicAttributes::name)
            .join();
   }
}
