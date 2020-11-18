package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public record WorldAttributes(String name, int flag, String message, String eventMessage, boolean recommended,
                              String recommendedMessage, int capacityStatus) implements AttributeResult {
}
