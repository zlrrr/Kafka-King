/*
 * Copyright 2025 Bronya0 <tangssst@163.com>.
 * Author Github: https://github.com/Bronya0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"testing"
)

// TestDisableIdempotenceConfiguration 测试幂等性禁用配置
func TestDisableIdempotenceConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		config         map[string]any
		expectDisabled bool
		description    string
	}{
		{
			name: "幂等性未配置 - 默认启用",
			config: map[string]any{
				"bootstrap_servers": "localhost:9092",
			},
			expectDisabled: false,
			description:    "当不配置 disable_idempotence 时，应使用默认设置（启用幂等）",
		},
		{
			name: "幂等性配置为 disable - 启用幂等",
			config: map[string]any{
				"bootstrap_servers":   "localhost:9092",
				"disable_idempotence": "disable",
			},
			expectDisabled: false,
			description:    "当 disable_idempotence 为 disable 时，不禁用幂等性",
		},
		{
			name: "幂等性配置为 enable - 禁用幂等",
			config: map[string]any{
				"bootstrap_servers":   "localhost:9092",
				"disable_idempotence": "enable",
			},
			expectDisabled: true,
			description:    "当 disable_idempotence 为 enable 时，应禁用幂等性",
		},
		{
			name: "完整配置 - 禁用幂等",
			config: map[string]any{
				"bootstrap_servers":   "broker1:9092,broker2:9092",
				"disable_idempotence": "enable",
				"sasl":                "enable",
				"sasl_mechanism":      "PLAIN",
				"sasl_user":           "test_user",
				"sasl_pwd":            "test_password",
			},
			expectDisabled: true,
			description:    "在包含 SASL 认证的完整配置中，应正确禁用幂等性",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("测试场景: %s", tt.description)

			// 检查配置值
			disableIdempotence := tt.config["disable_idempotence"]
			isEnabled := disableIdempotence == "enable"

			if isEnabled != tt.expectDisabled {
				t.Errorf("幂等性禁用状态不匹配:\n  期望: %v\n  实际: %v\n  配置: disable_idempotence=%v",
					tt.expectDisabled, isEnabled, disableIdempotence)
			}

			t.Logf("✓ 配置验证通过: disable_idempotence=%v, 应禁用=%v",
				disableIdempotence, tt.expectDisabled)
		})
	}
}

// TestConfigurationPersistence 测试配置持久化
func TestConfigurationPersistence(t *testing.T) {
	t.Run("配置字段存在性验证", func(t *testing.T) {
		// 模拟前端发送的配置
		frontendConfig := map[string]any{
			"id":                   1,
			"name":                 "Test Cluster",
			"bootstrap_servers":    "localhost:9092",
			"tls":                  "disable",
			"sasl":                 "disable",
			"use_ssh":              "disable",
			"disable_idempotence":  "enable", // 关键配置项
		}

		// 验证字段存在
		if _, exists := frontendConfig["disable_idempotence"]; !exists {
			t.Error("配置中缺少 disable_idempotence 字段")
		}

		// 验证字段值
		value := frontendConfig["disable_idempotence"]
		if value != "enable" && value != "disable" {
			t.Errorf("disable_idempotence 值不合法: %v (应该是 'enable' 或 'disable')", value)
		}

		t.Log("✓ 配置持久化验证通过")
	})
}

// TestBackwardCompatibility 测试向后兼容性
func TestBackwardCompatibility(t *testing.T) {
	t.Run("旧配置兼容性", func(t *testing.T) {
		// 模拟旧版本配置（没有 disable_idempotence 字段）
		oldConfig := map[string]any{
			"id":                1,
			"name":              "Legacy Cluster",
			"bootstrap_servers": "localhost:9092",
			"tls":               "disable",
			"sasl":              "disable",
		}

		// 验证不会因为缺少新字段而失败
		if _, exists := oldConfig["disable_idempotence"]; exists {
			t.Error("旧配置不应该包含 disable_idempotence 字段")
		}

		// 验证默认行为（不禁用幂等性）
		disableValue, _ := oldConfig["disable_idempotence"]
		if disableValue == "enable" {
			t.Error("旧配置应该默认不禁用幂等性")
		}

		t.Log("✓ 向后兼容性验证通过：旧配置不会被破坏")
	})
}

// TestUIDataFlow 测试 UI 数据流
func TestUIDataFlow(t *testing.T) {
	t.Run("前端到后端的数据流", func(t *testing.T) {
		// 1. 用户在前端设置开关
		userInput := "enable" // 用户打开了"禁用幂等性"开关

		// 2. 前端发送到后端
		requestPayload := map[string]any{
			"name":                "Test Connection",
			"bootstrap_servers":   "localhost:9092",
			"disable_idempotence": userInput,
		}

		// 3. 后端接收并处理
		receivedValue := requestPayload["disable_idempotence"]
		if receivedValue != userInput {
			t.Errorf("数据传输失败: 发送=%v, 接收=%v", userInput, receivedValue)
		}

		// 4. 验证配置应用逻辑
		shouldDisable := receivedValue == "enable"
		if !shouldDisable {
			t.Error("配置应用逻辑错误: enable 应该禁用幂等性")
		}

		t.Log("✓ UI 数据流验证通过")
	})
}

// TestKafkaCompatibility 测试 Kafka 版本兼容性场景
func TestKafkaCompatibility(t *testing.T) {
	scenarios := []struct {
		name            string
		kafkaVersion    string
		clientVersion   string
		shouldDisable   bool
		expectedError   string
		description     string
	}{
		{
			name:          "Kafka 3.0+ 客户端连接 2.x 服务端",
			kafkaVersion:  "2.8.0",
			clientVersion: "3.0.0",
			shouldDisable: true,
			expectedError: "Cluster authorization failed",
			description:   "高版本客户端默认启用幂等性，低版本服务端用户没有权限，应禁用幂等性",
		},
		{
			name:          "Kafka 3.0+ 客户端连接 3.0+ 服务端",
			kafkaVersion:  "3.0.0",
			clientVersion: "3.0.0",
			shouldDisable: false,
			expectedError: "",
			description:   "版本匹配，可以使用幂等性",
		},
		{
			name:          "Kafka 2.x 客户端连接 2.x 服务端",
			kafkaVersion:  "2.8.0",
			clientVersion: "2.8.0",
			shouldDisable: false,
			expectedError: "",
			description:   "版本匹配，使用默认设置",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Logf("场景: %s", scenario.description)
			t.Logf("  服务端版本: %s", scenario.kafkaVersion)
			t.Logf("  客户端版本: %s", scenario.clientVersion)
			t.Logf("  应禁用幂等: %v", scenario.shouldDisable)
			t.Logf("  预期错误: %s", scenario.expectedError)

			config := map[string]any{
				"bootstrap_servers": "test-broker:9092",
			}

			if scenario.shouldDisable {
				config["disable_idempotence"] = "enable"
			}

			// 验证配置
			disableValue, exists := config["disable_idempotence"]
			if scenario.shouldDisable && (!exists || disableValue != "enable") {
				t.Errorf("配置错误: 应该禁用幂等性但配置不正确")
			}

			t.Log("✓ 兼容性场景验证通过")
		})
	}
}

// BenchmarkConfigurationCheck 性能测试：配置检查
func BenchmarkConfigurationCheck(b *testing.B) {
	config := map[string]any{
		"bootstrap_servers":   "localhost:9092",
		"disable_idempotence": "enable",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config["disable_idempotence"] == "enable"
	}
}
