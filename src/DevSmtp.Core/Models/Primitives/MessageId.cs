namespace DevSmtp.Core.Models
{
    public class MessageId : IEquatable<MessageId>
    {
        public MessageId(string value)
        {
            this.Value = value;

            this.Validate();
        }

        public string Value { get; }

        public static MessageId From(string value) => new(value);

        public static MessageId NewId() => new(Guid.NewGuid().ToString().ToLower());

        private void Validate()
        {
            if (string.IsNullOrWhiteSpace(this.Value))
            {
                throw new FormatException($"{nameof(this.Value)} cannot be empty.");
            }
        }

        public bool Equals(MessageId? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return this.Value == other.Value;
        }

        public override bool Equals(object? obj)
        {
            return ReferenceEquals(this, obj) || obj is MessageId other && Equals(other);
        }

        public override int GetHashCode()
        {
            return this.Value.GetHashCode();
        }
    }
}
