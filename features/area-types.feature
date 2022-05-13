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
                "totalCount": 2,
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
        "area-types":[
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
        ],
        "limit": 20,
        "offset": 0,
        "count": 2,
        "total_count": 2
      }
      """

    And the HTTP status code should be "200"


  Scenario: Getting area-types no dataset

    When the following geography query response is available from Cantabular api extension for the dataset "":
      """
      {
        "data":{
          "dataset":{
            "ruleBase":{
              "isSourceOf":{
                "totalCount": 2,
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
                                "filterOnly": "false",
                                "label": "City",
                                "name": "city"
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
                      "label": "City",
                      "mapFrom": [],
                      "name": "city"
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
            "id": "country",
            "label": "Country",
            "total_count": 2
          },
          {
            "id": "city",
            "label": "City",
            "total_count": 3
          }
        ],
        "limit": 20,
        "offset": 0,
        "count": 2,
        "total_count": 2
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

    And I GET "/area-types?dataset=Inexistent"

    Then I should receive the following JSON response:
      """
      {
        "errors":[
          "failed to get area-types: error(s) returned by graphQL query"
        ]
      }
      """

    And the HTTP status code should be "404"

  Scenario: Getting area-types invalid query parameters

    When the following geography query response is available from Cantabular api extension for the dataset "Example":
      """
      {
        "data":{
          "dataset":{
            "ruleBase":{
              "isSourceOf":{
                "totalCount": 2,
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

    When I GET "/area-types?dataset=Example&limit=-1"
    
    Then I should receive the following JSON response:
      """
      {
        "errors":[
          "failed to parse request: invalid request: 'limit' cannot be a negative value"
        ]
      }
      """

    And the HTTP status code should be "400"

    When I GET "/area-types?dataset=Example&offset=-1"
    
    Then I should receive the following JSON response:
      """
      {
        "errors":[
          "failed to parse request: invalid request: 'offset' cannot be a negative value"
        ]
      }
      """

    And the HTTP status code should be "400"
