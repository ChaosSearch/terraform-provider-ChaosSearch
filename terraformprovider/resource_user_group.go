package main

import (
	"context"
	"cs-tf-provider/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_groups": {
				Type:     schema.TypeSet,
				Optional: false,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": {
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
							Optional: true,
							Computed: true,
						},

						"permissions": {
							Type:     schema.TypeSet,
							Optional: false,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"permission": {
										Type: schema.TypeSet,
										//Required: true,
										Optional: true,
										ForceNew: false,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"effect": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"action": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"resources": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"version": {
													Type:     schema.TypeString,
													Required: false,
													ForceNew: false,
													Optional: true,
													Computed: true,
												},
												"conditions": {
													Type:     schema.TypeSet,
													Optional: true,
													ForceNew: false,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"condition": {
																Type:     schema.TypeSet,
																Optional: false,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"starts_with": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"equals": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"not_equals": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"like": {
																			Type:     schema.TypeSet,
																			Optional: true,
																			ForceNew: false,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"chaos_document_attributes_title": {
																						Type:     schema.TypeString,
																						Required: false,
																						ForceNew: false,
																						Optional: true,
																						Computed: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("creating groups")
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	log.Debug("token value-->", tokenValue)

	var id, name string
	var permissionList []client.Permission
	var conditionList []client.Condition
	var actionsList []string
	var resourcesList []string

	if data.Get("user_groups").(*schema.Set).Len() > 0 {
		userGroupInterface := data.Get("user_groups").(*schema.Set).List()[0].(map[string]interface{})
		id = userGroupInterface["id"].(string)
		name = userGroupInterface["name"].(string)
		permissions := userGroupInterface["permissions"].(*schema.Set).List()[0].(map[string]interface{})["permission"].(*schema.Set).List()
		if len(permissions) > 0 {
			for permissionsIndex, permissionsElement := range permissions {

				permissionMap := permissionsElement.(map[string]interface{})
				var ConditionGroup client.ConditionGroup
				if len(permissionMap["conditions"].(*schema.Set).List()) > 0 {
					conditions := permissionMap["conditions"].(*schema.Set).List()[0].(map[string]interface{})["condition"].(*schema.Set).List()
					conditionMap := conditions[0].(map[string]interface{})
					equal := conditionMap["equals"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					startsWith := conditionMap["starts_with"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					notEquals := conditionMap["not_equals"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)
					like := conditionMap["like"].(*schema.Set).List()[0].(map[string]interface{})["chaos_document_attributes_title"].(string)

					equalObject := client.Equals{
						ChaosDocumentAttributesTitle: equal,
					}
					likeObject := client.Like{
						ChaosDocumentAttributesTitle: like,
					}
					notEqualsObject := client.NotEquals{
						ChaosDocumentAttributesTitle: notEquals,
					}
					startsWithObject := client.StartsWith{
						ChaosDocumentAttributesTitle: startsWith,
					}
					conditionList = append(
						conditionList,
						client.Condition{
							Equals:     equalObject,
							StartsWith: startsWithObject,
							NotEquals:  notEqualsObject,
							Like:       likeObject,
						})
					ConditionGroup = client.ConditionGroup{
						Condition: conditionList,
					}

				}

				actionsList = append(actionsList, permissionMap["action"].(string))
				resourcesList = append(resourcesList, permissionMap["resources"].(string))
				permissionList = append(
					permissionList,
					client.Permission{
						Effect:         permissionMap["effect"].(string),
						Actions:        actionsList,
						Resources:      resourcesList,
						Version:        permissionMap["version"].(string),
						ConditionGroup: ConditionGroup,
					})
				log.Debug(permissionsIndex, "index")
				//remove element from condition list after append
				conditionList = nil
			}
			log.Info("permission array", permissionList)
		}
	}

	createUserGroupRequest := &client.CreateUserGroupRequest{
		AuthToken:  tokenValue,
		Id:         id,
		Name:       name,
		Permission: permissionList,
	}

	if err := c.CreateUserGroup(ctx, createUserGroupRequest); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(name)

	return nil
}

func resourceGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))
	diags := diag.Diagnostics{}
	c := meta.(*ProviderMeta).CSClient
	tokenValue := meta.(*ProviderMeta).token
	req := &client.ReadUserGroupRequest{
		AuthToken: tokenValue,
		ID:        data.Id(),
	}
	resp, err := c.ReadUserGroup(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.Errorf("Couldn't find User Group: %s", err)
	}
	userGroupContent := make([]map[string]interface{}, 1)
	permissionContent := make(map[string]interface{})
	result := make([]map[string]interface{}, 1)
	userGroupContentMap := make(map[string]interface{})
	if resp.Permissions != nil && len(resp.Permissions) > 0 {
		permissions := make([]interface{}, len(resp.Permissions))
		for i := 0; i < len(resp.Permissions); i++ {
			permissionContent["effect"] = resp.Permissions[i].Effect
			permissionContent["actions"] = resp.Permissions[i].Actions
			permissionContent["resources"] = resp.Permissions[i].Resources
			permissionContent["version"] = resp.Permissions[i].Version
			permissions[i] = permissionContent
		}
		userGroupContentMap["permissions"] = permissions
	}
	userGroupContentMap["id"] = resp.Id
	userGroupContentMap["name"] = resp.Name
	result[0] = userGroupContentMap
	userGroupContent = result
	if err := data.Set("user_group", userGroupContent); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("updating groups")
	// c := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value-->", tokenValue)
	//TODO to be developed
	return nil
}

func resourceGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Debug("deleting groups")
	// c := meta.(*ProviderMeta).Client
	tokenValue := meta.(*ProviderMeta).token
	log.Warn("token value-->", tokenValue)
	//TODO to be developed
	return nil
}
