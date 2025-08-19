resource "stacklet_notification_template" "example" {
  name        = "example"
  description = "An email template"
  transport   = "email"
  content     = <<EOT
<!DOCTYPE html>
<html lang="en">
  <body>
    {% for m in messages %}
      {% set policy = m["policy"] %}
      {% set resources = m["resources"] %}
      {% set account = m["account"] %}
      {% set region = m["region"] %}
      {% set action = m["action"] %}
      
      <ul>
        <li><strong>Account:</strong> {{ account }}</li>
        <li><strong>Region:</strong> {{ region }}</li>
        <li><strong>Policy Name:</strong> {{ policy }}</li>
        <li><strong>Policy Description:</strong> {{ m["policy_description"] }}</li>
        <li><strong>Violation Description:</strong> {{ action['violation_desc'] }}</li>
        <li><strong>Resources:</strong></li>
        <li>{{ getResources(policy, resources) }}</li>
        <li><strong>Action Description:</strong> {{ action['action_desc'] }}</li>
      </ul>
    {% endfor %}
  </body>
</html>
EOT
}
