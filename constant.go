package main

const (
	ApiHost          string = "127.0.0.1:16443"
	SchedulerName    string = "mikischeduler"
	NodesEndpoint    string = "https://127.0.0.1:16443/api/v1/nodes"
	PodsEndpoint     string = "https://127.0.0.1:16443/api/v1/pods"
	AuthToken        string = "eyJhbGciOiJSUzI1NiIsImtpZCI6IndQU0tqWGRRdTh4MzI3NTA3UGplbTUxcnZhVVFVSHlVdEpzdTJiLUxhNmsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImFwaS1hY2Nlc3MtdG9rZW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiYXBpLWFjY2VzcyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjAyZDI1MmQ2LTBiYTQtNGMwOS04NzIzLTdjM2RmOTlhZGEwZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmFwaS1hY2Nlc3MifQ.mserqd-y0bzi1ZLdDXRwFXzdzFmufuHp6qHvQZ9jsH4OmAqX8Ns1yV2QwTAXeo_Ln1J-7OIU9JAEn-Ir-LqlQ5ZZihwXuXksYWNunfQhmpDc1mfqU5IU9kJviwZfoOwMjwdy_ciwbnVpY6nhoq5zzR4s5fUmEf-yr1dWylAV0Bght94bSskwDLYyIbtkoIu5oBzEXIvc-Xn7h8-ihr_SZtl_vvU7w4XeKma5wn_8c6RQvE1y7Ax9tl2HhicWcfpxuzk3A6lfSLytX2z70hxcronqECBdwP8WacMmhwlUo7DyTT1u9sJJhktG0bpqWOgdWDh5cMO5X0pErnJdc60GEQ"
	PodsEndpointv2   string = "/api/v1/pods"
	BindingsEndpoint string = "/api/v1/namespaces/default/pods/%s/binding/"
)
