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

package types

type Tag struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}
type Config struct {
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Language string    `json:"language"`
	Theme    string    `json:"theme"`
	Connects []Connect `json:"connects"`
}
type ResultsResp struct {
	Results []any  `json:"results"`
	Err     string `json:"err"`
}
type ResultResp struct {
	Result map[string]any `json:"result"`
	Err    string         `json:"err"`
}
type Connect struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	BootstrapServers    string `json:"bootstrap_servers"`
	Tls                 string `json:"tls"`
	SkipTLSVerify       string `json:"skipTLSVerify"`
	TlsCertFile         string `json:"tls_cert_file"`
	TlsKeyFile          string `json:"tls_key_file"`
	TlsCaFile           string `json:"tls_ca_file"`
	Sasl                string `json:"sasl"`
	SaslMechanism       string `json:"sasl_mechanism"`
	SaslUser            string `json:"sasl_user"`
	SaslPwd             string `json:"sasl_pwd"`
	KerberosUserKeytab  string `json:"kerberos_user_keytab"`
	KerberosKrb5Conf    string `json:"kerberos_krb5_conf"`
	KerberosUser        string `json:"Kerberos_user"`
	KerberosRealm       string `json:"Kerberos_realm"`
	KerberosServiceName string `json:"kerberos_service_name"`
	UseSsh              string `json:"use_ssh"`      // 新增：是否使用 SSH
	SshHost             string `json:"ssh_host"`     // SSH 主机
	SshPort             int    `json:"ssh_port"`     // SSH 端口
	SshUser             string `json:"ssh_user"`     // SSH 用户名
	SshPassword         string `json:"ssh_password"` // SSH 密码
	SshKeyFile          string `json:"ssh_key_file"` // SSH 私钥文件
	DisableIdempotence  string `json:"disable_idempotence"` // 是否禁用幂等性（解决版本兼容问题）
}
type H map[string]any

type GroupInfo struct {
	Group    string
	Topics   map[string][]PartitionOffset
	TotalLag int64
}

type PartitionOffset struct {
	Partition int32
	Offset    int64
	Lag       int64
}

// BrokerConfig Broker配置结构
type BrokerConfig struct {
	Name      string          `json:"name"`
	Value     string          `json:"value"`
	Source    string          `json:"source"` // DYNAMIC_BROKER_CONFIG, DEFAULT_CONFIG, STATIC_BROKER_CONFIG
	ReadOnly  bool            `json:"read_only"`
	Sensitive bool            `json:"sensitive"`
	Synonyms  []ConfigSynonym `json:"synonyms,omitempty"`
}

type ConfigSynonym struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}
