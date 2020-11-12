package com.atlas.wrg.processor;

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
}
