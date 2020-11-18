package com.atlas.wrg.rest.attribute;

import rest.AttributeResult;

public record WorldAttributes(String name, Integer flag, String message, String eventMessage, boolean recommended,
                              String recommendedMessage, Integer capacityStatus) implements AttributeResult {
}
