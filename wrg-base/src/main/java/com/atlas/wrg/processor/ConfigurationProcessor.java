package com.atlas.wrg.processor;

import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.List;
import java.util.Optional;

import com.atlas.wrg.configuration.Configuration;
import com.atlas.wrg.configuration.WorldConfiguration;
import com.esotericsoftware.yamlbeans.YamlReader;

public class ConfigurationProcessor {
   private static final Object lock = new Object();

   private static volatile ConfigurationProcessor instance;

   private final Configuration configuration;

   public static ConfigurationProcessor getInstance() {
      ConfigurationProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ConfigurationProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   private ConfigurationProcessor() {
      String fileName = "/service/config.yaml";
      String message;
      try {
         YamlReader reader = new YamlReader(new FileReader(fileName));
         configuration = reader.read(Configuration.class);
         reader.close();
      } catch (FileNotFoundException var3) {
         message = "Could not read config file " + fileName + ": " + var3.getMessage();
         throw new RuntimeException(message);
      } catch (IOException var4) {
         message = "Could not successfully parse config file " + fileName + ": " + var4.getMessage();
         throw new RuntimeException(message);
      }
   }

   public Optional<WorldConfiguration> getWorldConfiguration(int index) {
      List<WorldConfiguration> worldConfigurationList = getConfiguration().worlds;
      if (index < 0 || index >= worldConfigurationList.size()) {
         return Optional.empty();
      }
      return Optional.of(worldConfigurationList.get(index));
   }

   public Configuration getConfiguration() {
      return configuration;
   }
}
