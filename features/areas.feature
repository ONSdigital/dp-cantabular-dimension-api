Feature: Areas

  Background:

    Given private endpoints are not enabled

    And cantabular server is healthy

    And cantabular api extension is healthy

  Scenario: Getting areas happy
    When the following area query response is available from Cantabular api extension for the dataset "Example":
      """
     {
  "data": {
    "dataset": {
      "ruleBase": {
        "isSourceOf": {
          "search": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "E",
                          "label": "England",
                          "variable": {
                            "name": "country"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "N",
                          "label": "Northern Ireland",
                          "variable": {
                            "name": "country"
                          }
                        }
                      }
                    ],
                    "totalCount": 2
                  },
                  "mapFrom": [
                    {
                      "edges": [
                        {
                          "node": {
                            "label": "City",
                            "name": "city"
                          }
                        }
                      ]
                    }
                  ],
                  "name": "country"
                }
              },
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "0",
                          "label": "London",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast",
                          "variable": {
                            "name": "city"
                          }
                        }
                      }
                    ],
                    "totalCount": 3
                  },
                  "mapFrom": [],
                  "name": "city"
                }
              }
            ]
          }
        }
      }
    }
  }
}
      """
    And I GET "/areas?dataset=Example"

    Then I should receive the following JSON response:
      """
      {
    "areas": [
        {
            "id": "E",
            "label": "England",
            "area-type": "country"
        },
        {
            "id": "N",
            "label": "Northern Ireland",
            "area-type": "country"
        },
        {
            "id": "0",
            "label": "London",
            "area-type": "city"
        },
        {
            "id": "1",
            "label": "Liverpool",
            "area-type": "city"
        },
        {
            "id": "2",
            "label": "Belfast",
            "area-type": "city"
        }
    ]
}
      """

    And the HTTP status code should be "200"

  Scenario: Getting areas specific search
    When the following area query response is available from Cantabular api extension for the dataset "Example" and text "City":
    """
    {
  "data": {
    "dataset": {
      "ruleBase": {
        "isSourceOf": {
          "search": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "0",
                          "label": "London",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast",
                          "variable": {
                            "name": "city"
                          }
                        }
                      }
                    ],
                    "totalCount": 3
                  },
                  "mapFrom": [],
                  "name": "city"
                }
              }
            ]
          }
        }
      }
    }
  }
}
    """

    And I GET "/areas?dataset=Example&text=City"

    Then I should receive the following JSON response:
    """
    {
    "areas": [
        {
            "id": "0",
            "label": "London",
            "area-type": "city"
        },
        {
            "id": "1",
            "label": "Liverpool",
            "area-type": "city"
        },
        {
            "id": "2",
            "label": "Belfast",
            "area-type": "city"
        }
    ]
}
    """

    Scenario: Getting areas no dataset or search text
      When the following area query response is available from Cantabular api extension for the dataset "" and text "":
      """
      {
  "data": {
    "dataset": {
      "ruleBase": {
        "isSourceOf": {
          "search": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "E",
                          "label": "England",
                          "variable": {
                            "name": "country"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "N",
                          "label": "Northern Ireland",
                          "variable": {
                            "name": "country"
                          }
                        }
                      }
                    ],
                    "totalCount": 2
                  },
                  "mapFrom": [
                    {
                      "edges": [
                        {
                          "node": {
                            "label": "City",
                            "name": "city"
                          }
                        }
                      ]
                    }
                  ],
                  "name": "country"
                }
              },
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "0",
                          "label": "London",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool",
                          "variable": {
                            "name": "city"
                          }
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast",
                          "variable": {
                            "name": "city"
                          }
                        }
                      }
                    ],
                    "totalCount": 3
                  },
                  "mapFrom": [],
                  "name": "city"
                }
              }
            ]
          }
        }
      }
    }
  }
}
      """

      And I GET "/areas"

    Then I should receive the following JSON response:
    """
    {
    "areas": [
        {
            "id": "E",
            "label": "England",
            "area-type": "country"
        },
        {
            "id": "N",
            "label": "Northern Ireland",
            "area-type": "country"
        },
        {
            "id": "0",
            "label": "London",
            "area-type": "city"
        },
        {
            "id": "1",
            "label": "Liverpool",
            "area-type": "city"
        },
        {
            "id": "2",
            "label": "Belfast",
            "area-type": "city"
        }
    ]
}
    """

    Scenario: Getting areas invalid dataset
      When the following area query response is available from Cantabular api extension for the dataset "Test":
      """
      {
  "data": {
    "dataset": null
  },
  "errors": [
    {
      "message": "404 Not Found: dataset not loaded in this server",
      "locations": [
        {
          "line": 2,
          "column": 2
        }
      ],
      "path": [
        "dataset"
      ]
    }
  ]
}
      """
    And I GET "/areas?dataset=Test"

    Then I should receive the following JSON response:
    """
    {
    "errors": [
        "failed to get areas: error(s) returned by graphQL query"
    ]
}
    """

     Scenario: Get areas area does not exist
       When the following area query response is available from Cantabular api extension for the dataset "Example" and text "rio":
       """
       {
  "data": {
    "dataset": {
      "ruleBase": {
        "isSourceOf": {
          "search": {
            "edges": []
          }
        }
      }
    }
  }
}
       """
       And I GET "/areas?dataset=Example&text=rio"

    Then I should receive the following JSON response:
    """
    {
    "areas": null
    }
    """