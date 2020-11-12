package com.atlas.wrg.model;

public enum ServerStatus {
   NORMAL(0), HIGHLY_POPULATED(1), FULL(2);

   private final int value;

   ServerStatus(int value) {
      this.value = value;
   }

   public int getValue() {
      return value;
   }

   public static ServerStatus fromValue(int value) {
      for (ServerStatus op : ServerStatus.values()) {
         if (op.getValue() == value) {
            return op;
         }
      }
      return null;
   }
}
