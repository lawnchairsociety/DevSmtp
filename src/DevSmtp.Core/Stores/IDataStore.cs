using DevSmtp.Core.Models;

namespace DevSmtp.Core.Stores
{
    public interface IDataStore
    {
        Task StoreAsync(Message message, CancellationToken cancellationToken = default);

        Task<IEnumerable<Message>> GetAsync(CancellationToken cancellationToken = default);

        Task<Message?> FindByIdAsync(MessageId id, CancellationToken cancellationToken = default);

        Task<IEnumerable<Message>> FindByEmailAsync(Email email, CancellationToken cancellationToken = default);
    }
}
