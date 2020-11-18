package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public record WorldAttributes(String name, int flag, String message, String eventMessage, String recommended,
                              String recommendedMessage, int capacityStatus) implements AttributeResult {
}
