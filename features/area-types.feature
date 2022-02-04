Feature: Area Types

  Background:

    Given private endpoints are not enabled

    And cantabular server is healthy

    And cantabular api extension is healthy


  Scenario: Getting area-types happy

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
    And I GET "/area-types?cantabular_dataset=Example"

    Then I should receive the following JSON response:
      """
      {
        "area-types":[
          {
            "id":"country",
            "label":"Country"
          },
          {
            "id":"city",
            "label":"City"
          }
        ]
      }
      """

    And the HTTP status code should be "200"


  Scenario: Getting area-types no cantabular_dataset

    When the following geography query response is available from Cantabular api extension for the dataset "":
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

    And I GET "/area-types"

    Then I should receive the following JSON response:
      """
      {
        "area-types":[
          {
            "id":"country",
            "label":"Country"
          },
          {
            "id":"city",
            "label":"City"
          }
        ]
      }
      """

    And the HTTP status code should be "200"


  Scenario: Getting area-types not found

    When the following error json response is returned from Cantabular api extension for the dataset "Inexistent":
      """
      {
        "data":{
          "dataset":null
        },
        "errors":[
          {
            "message":"404 Not Found: dataset not loaded in this server",
            "locations":[
              {
                "line":2,
                "column":3
              }
            ],
            "path":[
              "dataset"
            ]
          }
        ]
      }
      """

    And I GET "/area-types?cantabular_dataset=Inexistent"

    Then I should receive the following JSON response:
      """
      {
        "errors":[
          "failed to get area-types from cantabular: error(s) returned by graphQL query"
        ]
      }
      """

    And the HTTP status code should be "404"
