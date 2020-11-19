package com.atlas.wrg.model;

import java.util.Arrays;
import java.util.Optional;

public enum ServerStatus {
   NORMAL(0), HIGHLY_POPULATED(1), FULL(2);

   private final int value;

   ServerStatus(int value) {
      this.value = value;
   }

   public int getValue() {
      return value;
   }

   public static Optional<ServerStatus> fromValue(int value) {
      return Arrays.stream(ServerStatus.values())
            .filter(serverStatus -> serverStatus.getValue() == value)
            .findAny();
   }
}
