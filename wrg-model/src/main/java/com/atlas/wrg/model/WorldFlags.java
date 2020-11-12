package com.atlas.wrg.model;

public enum WorldFlags {
   NOTHING(0),
   EVENT(1),
   NEW(2),
   HOT(3);

   private final int value;

   public int getValue() {
      return value;
   }

   WorldFlags(int value) {
      this.value = value;
   }
}
