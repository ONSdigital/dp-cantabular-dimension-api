Feature: Area Types Private Endpoints Enabled

  Background:

    Given private endpoints are enabled

    And cantabular server is healthy

    And cantabular api extension is healthy


  Scenario: Getting area-types when authorised

    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When the following geography query response is available from Cantabular api extension for the dataset "Example":
      """
      {
        "data":{
          "dataset":{
            "ruleBase":{
              "isSourceOf":{
                "edges":[
                  {
                    "node":{
                      "categories":{
                        "totalCount":2
                      },
                      "label":"Country",
                      "mapFrom":[
                        {
                          "edges":[
                            {
                              "node":{
                                "filterOnly":"false",
                                "label":"City",
                                "name":"city"
                              }
                            }
                          ]
                        }
                      ],
                      "name":"country"
                    }
                  },
                  {
                    "node":{
                      "categories":{
                        "totalCount":3
                      },
                      "label":"City",
                      "mapFrom":[],
                      "name":"city"
                    }
                  }
                ]
              },
              "name":"city"
            }
          }
        }
      }
      """

    And I GET "/area-types?dataset=Example"

    Then I should receive the following JSON response:
      """
      {
        "area-types": [
          {
            "id": "country",
            "label": "Country",
            "total_count": 2
          },
          {
            "id": "city",
            "label": "City",
            "total_count": 3
          }
        ]
      }
      """

    And the HTTP status code should be "200"


  Scenario: Getting area-types when not authorised

    Given I am not identified

    And I am not authorised

    When I GET "/area-types?dataset=Example"

    Then the HTTP status code should be "401"
